// 离线将 sys_member.password 中历史 bcrypt 批量改为 SM2 密文（与运行时 EncryptMemberPassword 一致）。
//
// 重要：bcrypt 不可逆，无法仅凭数据库哈希自动得到明文。必须提供含明文的 CSV，
// 脚本会校验「CSV 明文」与「库中 bcrypt」匹配后再写入 SM2。
//
// 用法（在项目根目录，且当前目录有可读 config.yaml，或复制 manifest/config/config.yaml 为 config.yaml）：
//
//	# 仅统计仍为 bcrypt 的会员数量与 id 列表
//	go run ./cmd/migrate_sys_member_password_sm2 -list
//
//	# 根据 CSV 批量迁移（CSV 列：member_id,plain_password）
//	go run ./cmd/migrate_sys_member_password_sm2 -csv ./member_plain.csv
//
//	# Excel / 记事本另存为「制表符分隔」时请加 -tsv（否则整行会被当成一列，表头解析失败）
//	go run ./cmd/migrate_sys_member_password_sm2 -csv ./demo.csv -tsv
//
//	# 先演练不写库
//	go run ./cmd/migrate_sys_member_password_sm2 -csv ./member_plain.csv -dry-run
//
//	# 把 CSV 里的明文批量转成 SM2 密文（hex）写入新文件，不写数据库（无 bcrypt 校验）
//	go run ./cmd/migrate_sys_member_password_sm2 -csv ./plain.csv -tsv -encrypt-file-only -sm2-out ./sm2_passwords.csv
//
//	# 迁移写库的同时，把每行生成的 SM2 密文另存一份（便于核对）
//	go run ./cmd/migrate_sys_member_password_sm2 -csv ./plain.csv -sm2-out ./sm2_passwords.csv
//
//	# 不校验旧密码，直接按 CSV 明文加密并更新数据库
//	go run ./cmd/migrate_sys_member_password_sm2 -csv ./plain.csv -tsv -no-check
//
// CSV 格式（UTF-8，首行为表头，逗号分隔）：
//
//	member_id,plain_password
//	1001,MySecret1
//	1002,AnotherPwd
//
// 也可用 username 列（二选一，与 member_id 不要混在同一文件除非每行只填一种）：
//
//	username,plain_password
//	user@example.com,MySecret1
package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"exam/internal/cmd"
	"exam/internal/consts"
	"exam/internal/dao"
	secsvc "exam/internal/service/security"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	_ "exam/internal/logic/security"
)

func main() {
	var (
		csvPath         = flag.String("csv", "", "含明文的 CSV 路径（列见文件头注释）")
		tsv             = flag.Bool("tsv", false, "使用制表符(Tab)作为列分隔符（Excel 导出常用）")
		dryRun          = flag.Bool("dry-run", false, "只打印将执行的更新，不写数据库")
		sm2OutPath      = flag.String("sm2-out", "", "将每行明文加密后的 SM2 密文(hex)写入此 CSV（列: member_id,username,password_sm2）")
		encryptFileOnly = flag.Bool("encrypt-file-only", false, "仅明文→SM2 写出 -sm2-out，不连库、不更新 sys_member（需 -sm2-out）")
		noCheck         = flag.Bool("no-check", false, "不读取旧密码也不做 bcrypt 校验，直接按 CSV 明文加密后更新 sys_member.password")
		listOnly        = flag.Bool("list", false, "仅列出仍为 bcrypt 的 sys_member id（不写库）")
		limit           = flag.Int("limit", 0, "与 -list 联用：最多输出多少条 id（0 表示不限制）")
	)
	flag.Parse()

	ctx := gctx.GetInitCtx()
	cmd.InitAll(ctx)

	if *listOnly {
		if err := runList(ctx, *limit); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if strings.TrimSpace(*csvPath) == "" {
		fmt.Fprintln(os.Stderr, "请指定 -csv 文件路径，或使用 -list 查看仍为 bcrypt 的记录。")
		fmt.Fprintln(os.Stderr, "说明：bcrypt 不可逆，没有明文无法自动转成 SM2。")
		os.Exit(2)
	}
	if *encryptFileOnly {
		if strings.TrimSpace(*sm2OutPath) == "" {
			fmt.Fprintln(os.Stderr, "-encrypt-file-only 必须同时指定 -sm2-out 输出路径")
			os.Exit(2)
		}
		if err := runEncryptPlainToSM2File(ctx, *csvPath, *tsv, *sm2OutPath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if err := runMigrate(ctx, *csvPath, *dryRun, *tsv, *sm2OutPath, *noCheck); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isBcryptStored(p string) bool {
	p = strings.TrimSpace(p)
	return strings.HasPrefix(p, "$2a$") || strings.HasPrefix(p, "$2b$") || strings.HasPrefix(p, "$2y$")
}

func runList(ctx context.Context, maxIDs int) error {
	model := dao.SysMember.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Wheref("(password LIKE ? OR password LIKE ? OR password LIKE ?)", "$2a$%", "$2b$%", "$2y$%")
	cnt, err := model.Count()
	if err != nil {
		return err
	}
	fmt.Printf("仍为 bcrypt 的 sys_member 行数: %d\n", cnt)
	if cnt == 0 {
		return nil
	}
	var rows []struct {
		Id int64 `json:"id"`
	}
	q := dao.SysMember.Ctx(ctx).
		Fields("id").
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Wheref("(password LIKE ? OR password LIKE ? OR password LIKE ?)", "$2a$%", "$2b$%", "$2y$%").
		OrderAsc("id")
	if maxIDs > 0 {
		q = q.Limit(maxIDs)
	}
	if err := q.Scan(&rows); err != nil {
		return err
	}
	fmt.Println("id 列表:")
	for _, r := range rows {
		fmt.Println(r.Id)
	}
	if maxIDs > 0 && len(rows) < cnt {
		fmt.Fprintf(os.Stderr, "提示: 已用 -limit=%d 截断输出，总行数仍为 %d\n", maxIDs, cnt)
	}
	return nil
}

func runMigrate(ctx context.Context, path string, dryRun, tsv bool, sm2OutPath string, noCheck bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	comma, err := csvDelimiterForRead(f, tsv)
	if err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	r := csv.NewReader(f)
	r.Comma = comma
	r.TrimLeadingSpace = true
	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("read csv header: %w", err)
	}
	header = normalizeCSVHeaderFields(header)
	colID, colUser, colPlain := resolveColumns(header)
	if colPlain < 0 || (colID < 0 && colUser < 0) {
		hint := ""
		if len(header) == 1 && strings.Contains(header[0], "\t") {
			hint = "（检测到表头在同一列内且含 Tab，请使用 -tsv 参数，或改为英文逗号分隔的 CSV）"
		}
		return fmt.Errorf("CSV 表头需包含 plain_password 以及 member_id 或 username 之一，当前: %v%s", header, hint)
	}

	var sm2Out *os.File
	var sm2W *csv.Writer
	if strings.TrimSpace(sm2OutPath) != "" {
		sm2Out, err = os.Create(sm2OutPath)
		if err != nil {
			return fmt.Errorf("create sm2-out: %w", err)
		}
		defer sm2Out.Close()
		sm2W = csv.NewWriter(sm2Out)
		sm2W.Comma = comma
		if err := sm2W.Write([]string{"member_id", "username", "password_sm2"}); err != nil {
			return err
		}
		defer func() {
			sm2W.Flush()
		}()
	}

	sec := secsvc.Security()
	var ok, skip, fail int
	line := 1
	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("csv line %d: %w", line+1, err)
		}
		line++
		if len(rec) == 0 || (len(rec) == 1 && strings.TrimSpace(rec[0]) == "") {
			continue
		}

		plain := strings.TrimSpace(safeCol(rec, colPlain))
		if plain == "" {
			fmt.Fprintf(os.Stderr, "skip line %d: empty plain_password\n", line)
			skip++
			continue
		}

		var memberID int64
		var csvUsername string
		if colID >= 0 {
			idStr := strings.TrimSpace(safeCol(rec, colID))
			var perr error
			memberID, perr = strconv.ParseInt(idStr, 10, 64)
			if perr != nil || memberID <= 0 {
				fmt.Fprintf(os.Stderr, "skip line %d: bad member_id %q\n", line, idStr)
				fail++
				continue
			}
			if colUser >= 0 {
				csvUsername = strings.TrimSpace(safeCol(rec, colUser))
			}
		} else {
			u := strings.TrimSpace(safeCol(rec, colUser))
			csvUsername = u
			if u == "" {
				fmt.Fprintf(os.Stderr, "skip line %d: empty username\n", line)
				skip++
				continue
			}
			var row struct {
				Id int64
			}
			qerr := dao.SysMember.Ctx(ctx).
				Fields("id").
				Where("username", u).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&row)
			if qerr != nil || row.Id == 0 {
				fmt.Fprintf(os.Stderr, "skip line %d: user not found username=%q\n", line, u)
				fail++
				continue
			}
			memberID = row.Id
		}

		if sm2W != nil && strings.TrimSpace(csvUsername) == "" && memberID > 0 {
			var nm struct {
				Username string
			}
			_ = dao.SysMember.Ctx(ctx).
				Fields("username").
				Where("id", memberID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&nm)
			csvUsername = strings.TrimSpace(nm.Username)
		}

		if !noCheck {
			var stored string
			err = dao.SysMember.Ctx(ctx).
				Fields("password").
				Where("id", memberID).
				Where("delete_flag", consts.DeleteFlagNotDeleted).
				Scan(&stored)
			if err != nil || strings.TrimSpace(stored) == "" {
				fmt.Fprintf(os.Stderr, "skip line %d: member id=%d load password failed\n", line, memberID)
				fail++
				continue
			}
			if !isBcryptStored(stored) {
				fmt.Fprintf(os.Stderr, "skip line %d: id=%d password is not bcrypt (already migrated?)\n", line, memberID)
				skip++
				continue
			}
			if !utility.CheckPassword(stored, plain) {
				fmt.Fprintf(os.Stderr, "skip line %d: id=%d plain does not match stored bcrypt\n", line, memberID)
				fail++
				continue
			}
		}

		sm2Hex, err := sec.EncryptMemberPassword(ctx, plain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line %d: id=%d encrypt failed: %v\n", line, memberID, err)
			fail++
			continue
		}

		if sm2W != nil {
			_ = sm2W.Write([]string{
				strconv.FormatInt(memberID, 10),
				csvUsername,
				sm2Hex,
			})
		}

		if dryRun {
			fmt.Printf("[dry-run] id=%d -> sm2 hex len=%d\n", memberID, len(sm2Hex))
			ok++
			continue
		}
		if _, err := dao.SysMember.Ctx(ctx).Where("id", memberID).Data(g.Map{
			"password": sm2Hex,
		}).Update(); err != nil {
			fmt.Fprintf(os.Stderr, "line %d: id=%d update failed: %v\n", line, memberID, err)
			fail++
			continue
		}
		fmt.Printf("ok id=%d\n", memberID)
		ok++
	}

	fmt.Fprintf(os.Stderr, "done: ok=%d skip=%d fail=%d\n", ok, skip, fail)
	return nil
}

// runEncryptPlainToSM2File 仅将 CSV 中 plain_password 转为 SM2 hex 写入文件，不访问数据库。
func runEncryptPlainToSM2File(ctx context.Context, path string, tsv bool, outPath string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	comma, err := csvDelimiterForRead(f, tsv)
	if err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	r := csv.NewReader(f)
	r.Comma = comma
	r.TrimLeadingSpace = true
	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("read csv header: %w", err)
	}
	header = normalizeCSVHeaderFields(header)
	colID, colUser, colPlain := resolveColumns(header)
	if colPlain < 0 || (colID < 0 && colUser < 0) {
		return fmt.Errorf("CSV 表头需包含 plain_password 以及 member_id 或 username 之一，当前: %v", header)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()
	w := csv.NewWriter(out)
	w.Comma = comma
	if err := w.Write([]string{"member_id", "username", "password_sm2"}); err != nil {
		return err
	}
	defer w.Flush()

	sec := secsvc.Security()
	var ok, fail int
	line := 1
	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("csv line %d: %w", line+1, err)
		}
		line++
		if len(rec) == 0 || (len(rec) == 1 && strings.TrimSpace(rec[0]) == "") {
			continue
		}
		plain := strings.TrimSpace(safeCol(rec, colPlain))
		if plain == "" {
			continue
		}
		sm2Hex, err := sec.EncryptMemberPassword(ctx, plain)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line %d: encrypt failed: %v\n", line, err)
			fail++
			continue
		}
		var idStr, userStr string
		if colID >= 0 {
			idStr = strings.TrimSpace(safeCol(rec, colID))
		}
		if colUser >= 0 {
			userStr = strings.TrimSpace(safeCol(rec, colUser))
		}
		if err := w.Write([]string{idStr, userStr, sm2Hex}); err != nil {
			return err
		}
		ok++
	}
	fmt.Fprintf(os.Stderr, "sm2-out 写入完成: %s （行数=%d, 失败=%d）\n", outPath, ok, fail)
	return nil
}

func csvDelimiterForRead(f *os.File, tsv bool) (rune, error) {
	if tsv {
		return '\t', nil
	}
	comma, err := sniffCSVDelimiter(f)
	if err != nil {
		return ',', err
	}
	return comma, nil
}

// sniffCSVDelimiter 根据首行判断常见「Tab 分隔」导出；默认逗号。
func sniffCSVDelimiter(f *os.File) (rune, error) {
	_, err := f.Seek(0, 0)
	if err != nil {
		return ',', err
	}
	br := bufio.NewReader(f)
	line, err := br.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return ',', err
	}
	line = strings.TrimPrefix(strings.TrimSpace(line), "\ufeff")
	tabN := strings.Count(line, "\t")
	commaN := strings.Count(line, ",")
	if tabN >= 1 && commaN == 0 {
		return '\t', nil
	}
	return ',', nil
}

func normalizeCSVHeaderFields(header []string) []string {
	out := make([]string, len(header))
	for i, h := range header {
		s := strings.TrimSpace(h)
		if i == 0 {
			s = strings.TrimPrefix(s, "\ufeff")
		}
		out[i] = s
	}
	return out
}

func resolveColumns(header []string) (idIdx, userIdx, plainIdx int) {
	idIdx, userIdx, plainIdx = -1, -1, -1
	for i, h := range header {
		k := strings.ToLower(strings.TrimSpace(h))
		switch k {
		case "member_id", "id":
			idIdx = i
		case "username", "user", "login", "email":
			if userIdx < 0 {
				userIdx = i
			}
		case "plain_password", "password", "plain", "pwd":
			plainIdx = i
		}
	}
	return idIdx, userIdx, plainIdx
}

func safeCol(rec []string, i int) string {
	if i < 0 || i >= len(rec) {
		return ""
	}
	return rec[i]
}
