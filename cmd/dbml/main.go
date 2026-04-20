// cmd/dbml/main.go
//
// Generates docs/schema.dbml from db/migrations/*.up.sql files.
// Publish to dbdocs.io via: dbdocs build docs/schema.dbml
//
// Usage:
//   go run ./cmd/dbml           → generate docs/schema.dbml
//   go run ./cmd/dbml -print    → dump DBML to stdout only

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ── flags ─────────────────────────────────────────────────────────────────────

var (
	flagPrint = flag.Bool("print", false, "Print DBML to stdout instead of writing file")
	flagDir   = flag.String("migrations", "db/migrations", "Path to migrations directory")
	flagOut   = flag.String("out", "docs/schema.dbml", "Output DBML file path")
)

type Column struct {
	Name        string
	Type        string
	Constraints []string // pk, not null, unique, default, note
}

type Index struct {
	Columns []string
	Unique  bool
	Note    string
}

type Table struct {
	Name    string
	Columns []Column
	Indexes []Index
	Note    string
}

type Ref struct {
	From       string // table.column
	To         string // table.column
	DeleteRule string // cascade | set null | restrict
	Relation   string // > (many-to-one) | - (one-to-one)
}

type Schema struct {
	Tables []Table
	Refs   []Ref
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("[dbml] ")

	// Read all *.up.sql files sorted by name
	pattern := filepath.Join(*flagDir, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		log.Fatalf("no *.up.sql files found in %s", *flagDir)
	}
	sort.Strings(files)

	var fullSQL strings.Builder
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}
		fullSQL.WriteString(string(data))
		fullSQL.WriteString("\n")
	}

	schema := parseSQL(fullSQL.String())
	dbml := renderDBML(schema)

	if *flagPrint {
		fmt.Println(dbml)
		return
	}

	// Write file
	if err := os.MkdirAll(filepath.Dir(*flagOut), 0755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(*flagOut, []byte(dbml), 0644); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Printf("wrote %s (%d tables, %d refs)", *flagOut, len(schema.Tables), len(schema.Refs))
	log.Printf("publish to dbdocs.io:\n  dbdocs build %s --project kreatip", *flagOut)
}

// ── SQL parser ────────────────────────────────────────────────────────────────

var (
	reCreateTable   = regexp.MustCompile(`(?is)CREATE\s+TABLE\s+(\w+)\s*\((.+?)\)\s*;`)
	reColumnLine    = regexp.MustCompile(`(?i)^\s*(\w+)\s+([\w()]+(?:\(\d+(?:,\d+)?\))?)(.*)$`)
	reDefaultVal    = regexp.MustCompile(`(?i)DEFAULT\s+('[^']*'|[^\s,]+)`)
	reCheckEnum     = regexp.MustCompile(`(?i)CONSTRAINT\s+\w+\s+CHECK\s*\((\w+)\s+IN\s*\(([^)]+)\)\)`)
	reUniqueConstr  = regexp.MustCompile(`(?i)CONSTRAINT\s+\w+\s+UNIQUE\s*\(([^)]+)\)`)
	reFKConstr      = regexp.MustCompile(`(?i)CONSTRAINT\s+\w+\s+FOREIGN\s+KEY\s*\(([^)]+)\)\s+REFERENCES\s+(\w+)\s*\(([^)]+)\)(?:\s+ON\s+DELETE\s+(\w[\w\s]*))?`)
	reCreateIndex   = regexp.MustCompile(`(?i)CREATE\s+(UNIQUE\s+)?INDEX\s+\w+\s+ON\s+(\w+)\s*\(([^)]+)\)(?:\s+WHERE\s+[^;]+)?`)
	reComment       = regexp.MustCompile(`--[^\n]*`)
	reBlockComment  = regexp.MustCompile(`(?s)/\*.*?\*/`)
	rePrimaryConstr = regexp.MustCompile(`(?i)CONSTRAINT\s+\w+\s+PRIMARY\s+KEY\s*\(([^)]+)\)`)
	reUniqueIdx     = regexp.MustCompile(`(?i)CONSTRAINT\s+\w+\s+UNIQUE\s*\(([^)]+)\)`)
)

func parseSQL(sql string) Schema {
	// Strip block comments, line comments
	sql = reBlockComment.ReplaceAllString(sql, " ")
	sql = reComment.ReplaceAllString(sql, "")

	schema := Schema{}
	tableMap := map[string]*Table{}

	for _, m := range reCreateTable.FindAllStringSubmatch(sql, -1) {
		tableName := m[1]
		body := m[2]

		table := &Table{Name: tableName}
		lines := splitColumnLines(body)

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			upper := strings.ToUpper(line)

			// PRIMARY KEY constraint
			if pkm := rePrimaryConstr.FindStringSubmatch(line); len(pkm) > 0 {
				for _, col := range splitCols(pkm[1]) {
					setConstraint(&table.Columns, col, "pk")
				}
				continue
			}

			// UNIQUE constraint
			if um := reUniqueConstr.FindStringSubmatch(line); len(um) > 0 {
				cols := splitCols(um[1])
				if len(cols) == 1 {
					setConstraint(&table.Columns, cols[0], "unique")
				} else {
					table.Indexes = append(table.Indexes, Index{Columns: cols, Unique: true})
				}
				continue
			}

			// FOREIGN KEY constraint
			if fkm := reFKConstr.FindStringSubmatch(line); len(fkm) > 0 {
				fromCol := strings.TrimSpace(fkm[1])
				toTable := strings.TrimSpace(fkm[2])
				toCol := strings.TrimSpace(fkm[3])
				deleteRule := strings.TrimSpace(strings.ToLower(fkm[4]))

				relation := ">"
				// Check if from col is unique → one-to-one
				if colIsUnique(table.Columns, fromCol) {
					relation = "-"
				}

				schema.Refs = append(schema.Refs, Ref{
					From:       tableName + "." + fromCol,
					To:         toTable + "." + toCol,
					DeleteRule: deleteRule,
					Relation:   relation,
				})
				continue
			}

			// CHECK constraint — skip (already handled via column notes)
			if strings.Contains(upper, "CONSTRAINT") && strings.Contains(upper, "CHECK") {
				continue
			}

			// Generic CONSTRAINT line
			if strings.HasPrefix(upper, "CONSTRAINT") {
				continue
			}

			if cm := reColumnLine.FindStringSubmatch(line); len(cm) > 0 {
				colName := cm[1]
				colType := normalizeType(cm[2])
				rest := cm[3]
				restUp := strings.ToUpper(rest)

				col := Column{Name: colName, Type: colType}

				// primary key inline
				if strings.Contains(restUp, "PRIMARY KEY") {
					col.Constraints = append(col.Constraints, "pk")
				}
				// not null
				if strings.Contains(restUp, "NOT NULL") {
					col.Constraints = append(col.Constraints, "not null")
				}
				// unique inline
				if strings.Contains(restUp, "UNIQUE") {
					col.Constraints = append(col.Constraints, "unique")
				}
				// default
				if dm := reDefaultVal.FindStringSubmatch(rest); len(dm) > 0 {
					val := dm[1]
					// wrap function calls in backticks, strings in single quotes
					if strings.HasPrefix(val, "'") {
						col.Constraints = append(col.Constraints, "default: "+val)
					} else {
						col.Constraints = append(col.Constraints, "default: `"+val+"`")
					}
				}

				table.Columns = append(table.Columns, col)
			}
		}

		tableMap[tableName] = table
		schema.Tables = append(schema.Tables, *table) // placeholder, will update below
	}

	// Rebuild tables slice from map (preserve order)
	tableOrder := []string{}
	for _, m := range reCreateTable.FindAllStringSubmatch(sql, -1) {
		tableOrder = append(tableOrder, m[1])
	}
	schema.Tables = nil
	for _, name := range tableOrder {
		if t, ok := tableMap[name]; ok {
			schema.Tables = append(schema.Tables, *t)
		}
	}

	for _, m := range reCreateIndex.FindAllStringSubmatch(sql, -1) {
		unique := strings.TrimSpace(m[1]) != ""
		tableName := m[2]
		cols := splitCols(m[3])

		for i := range schema.Tables {
			if schema.Tables[i].Name == tableName {
				// skip if single col already marked unique
				if len(cols) == 1 && colIsUnique(schema.Tables[i].Columns, cols[0]) {
					break
				}
				schema.Tables[i].Indexes = append(schema.Tables[i].Indexes, Index{
					Columns: cols,
					Unique:  unique,
				})
				break
			}
		}
	}

	for i, ref := range schema.Refs {
		parts := strings.SplitN(ref.From, ".", 2)
		if len(parts) != 2 {
			continue
		}
		tName, cName := parts[0], parts[1]
		for _, t := range schema.Tables {
			if t.Name == tName && colIsUnique(t.Columns, cName) {
				schema.Refs[i].Relation = "-"
			}
		}
	}

	return schema
}

func renderDBML(schema Schema) string {
	var b strings.Builder

	b.WriteString("// Kreatip — Database Schema\n")
	b.WriteString("// Auto-generated by cmd/dbml — DO NOT EDIT MANUALLY\n")
	b.WriteString("// Run: go run ./cmd/dbml -open\n\n")

	b.WriteString("Project kreatip {\n")
	b.WriteString("  database_type: 'PostgreSQL'\n")
	b.WriteString("  Note: 'Donation platform for creators & streamers'\n")
	b.WriteString("}\n\n")

	for _, table := range schema.Tables {
		b.WriteString(fmt.Sprintf("Table %s {\n", table.Name))

		for _, col := range table.Columns {
			parts := []string{col.Name, col.Type}
			if len(col.Constraints) > 0 {
				parts = append(parts, "["+strings.Join(col.Constraints, ", ")+"]")
			}
			b.WriteString("  " + strings.Join(parts, " ") + "\n")
		}

		if len(table.Indexes) > 0 {
			b.WriteString("\n  indexes {\n")
			for _, idx := range table.Indexes {
				// Strip ASC/DESC from each column — not valid in DBML index syntax
				clean := make([]string, len(idx.Columns))
				for i, c := range idx.Columns {
					c = strings.TrimSuffix(strings.TrimSpace(c), " DESC")
					c = strings.TrimSuffix(c, " ASC")
					c = strings.TrimSuffix(strings.TrimSpace(c), " desc")
					c = strings.TrimSuffix(c, " asc")
					clean[i] = strings.TrimSpace(c)
				}
				colStr := ""
				if len(clean) == 1 {
					colStr = clean[0]
				} else {
					colStr = "(" + strings.Join(clean, ", ") + ")"
				}
				if idx.Unique {
					b.WriteString(fmt.Sprintf("    %s [unique]\n", colStr))
				} else {
					b.WriteString(fmt.Sprintf("    %s\n", colStr))
				}
			}
			b.WriteString("  }\n")
		}

		b.WriteString("}\n\n")
	}

	for _, ref := range schema.Refs {
		line := fmt.Sprintf("Ref: %s %s %s", ref.From, ref.Relation, ref.To)
		if ref.DeleteRule != "" && ref.DeleteRule != "restrict" {
			line += fmt.Sprintf(" [delete: %s]", ref.DeleteRule)
		}
		b.WriteString(line + "\n")
	}

	return b.String()
}

// ── helpers ───────────────────────────────────────────────────────────────────

// splitColumnLines splits CREATE TABLE body by commas but respects parentheses depth.
func splitColumnLines(body string) []string {
	var lines []string
	depth := 0
	var cur strings.Builder
	for _, r := range body {
		switch r {
		case '(':
			depth++
			cur.WriteRune(r)
		case ')':
			depth--
			cur.WriteRune(r)
		case ',':
			if depth == 0 {
				lines = append(lines, cur.String())
				cur.Reset()
			} else {
				cur.WriteRune(r)
			}
		default:
			cur.WriteRune(r)
		}
	}
	if s := strings.TrimSpace(cur.String()); s != "" {
		lines = append(lines, s)
	}
	return lines
}

func splitCols(s string) []string {
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		// Strip sort direction — DBML indexes don't support ASC/DESC
		col := strings.TrimSpace(p)
		col = strings.TrimSuffix(col, " DESC")
		col = strings.TrimSuffix(col, " ASC")
		col = strings.TrimSuffix(col, " desc")
		col = strings.TrimSuffix(col, " asc")
		out = append(out, col)
	}
	return out
}

func setConstraint(cols *[]Column, colName string, constraint string) {
	for i, c := range *cols {
		if c.Name == colName {
			(*cols)[i].Constraints = append((*cols)[i].Constraints, constraint)
			return
		}
	}
}

func colIsUnique(cols []Column, name string) bool {
	for _, c := range cols {
		if c.Name == name {
			for _, ct := range c.Constraints {
				if ct == "unique" || ct == "pk" {
					return true
				}
			}
		}
	}
	return false
}

// normalizeType maps PG types to DBML-friendly names.
func normalizeType(t string) string {
	up := strings.ToUpper(t)
	switch {
	case up == "TIMESTAMPTZ":
		return "timestamp"
	case strings.HasPrefix(up, "VARCHAR"):
		return strings.ToLower(t)
	case up == "TEXT":
		return "text"
	case up == "BOOLEAN" || up == "BOOL":
		return "boolean"
	case up == "BIGINT":
		return "bigint"
	case up == "INT" || up == "INTEGER":
		return "integer"
	case up == "UUID":
		return "uuid"
	case up == "JSONB":
		return "jsonb"
	default:
		return strings.ToLower(t)
	}
}
