package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/hiidy/simpleblog/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const helpText = `Usage: main [flags] arg [arg...]

This is a pflag example.

Flags:
`

type Querier interface {
	FilterWithNameAndRole(name string) ([]gen.T, error)
}

type GenerateConfig struct {
	ModelPackagePath string
	GenerateFunc     func(g *gen.Generator)
}

var generateConfigs = map[string]GenerateConfig{
	"sb": {ModelPackagePath: "../../internal/apiserver/model", GenerateFunc: GenerateSimpleBlogModels},
}

var (
	addr       = pflag.StringP("addr", "a", "127.0.0.1:3306", "MySQL host address.")
	username   = pflag.StringP("username", "u", "simpleblog", "Username to connect to the database.")
	password   = pflag.StringP("password", "p", "simpleblog1234", "Password to use when connecting to the database.")
	database   = pflag.StringP("db", "d", "simpleblog", "Database name to connect to.")
	modelPath  = pflag.String("model-pkg-path", "", "Generated model code's package name.")
	components = pflag.StringSlice("component", []string{"sb"}, "Generated model code's for specified component.")
	help       = pflag.BoolP("help", "h", false, "Show this help message.")
)

func main() {
	pflag.Usage = func() {
		fmt.Printf("%s", helpText)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}

	dbInstance, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	for _, component := range *components {
		processComponent(component, dbInstance)
	}
}

func initializeDatabase() (*gorm.DB, error) {
	dbOptions := &db.MySQLOptions{
		Addr:     *addr,
		Username: *username,
		Password: *password,
		Database: *database,
	}

	return db.NewMySQL(dbOptions)
}

func processComponent(component string, dbInstance *gorm.DB) {
	config, ok := generateConfigs[component]
	if !ok {
		log.Printf("Component '%s' not found in configuration. Skipping.", component)
		return
	}

	modelPkgPath := resolveModelPackagePath(config.ModelPackagePath)

	generator := createGenerator(modelPkgPath)
	generator.UseDB(dbInstance)

	applyGeneratorOptions(generator)

	config.GenerateFunc(generator)

	generator.Execute()
}

func resolveModelPackagePath(defaultPath string) string {
	if *modelPath != "" {
		return *modelPath
	}
	absPath, err := filepath.Abs(defaultPath)
	if err != nil {
		log.Printf("Error resolving path: %v", err)
		return defaultPath
	}
	return absPath
}

func createGenerator(packagePath string) *gen.Generator {
	return gen.NewGenerator(gen.Config{
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		ModelPkgPath:  packagePath,
		WithUnitTest:  true,
		FieldNullable: true,
		FieldSignable: false, FieldWithIndexTag: false, FieldWithTypeTag: false,
	})
}

func applyGeneratorOptions(g *gen.Generator) {
	g.WithOpts(
		gen.FieldGORMTag("createdAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
		gen.FieldGORMTag("updatedAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
	)
}

func GenerateSimpleBlogModels(g *gen.Generator) {
	g.GenerateModelAs(
		"user",
		"UserM",
		gen.FieldIgnore("placeholder"),
		gen.FieldGORMTag("username", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_username")
			return tag
		}),
		gen.FieldGORMTag("userID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_userID")
			return tag
		}),
		gen.FieldGORMTag("phone", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_phone")
			return tag
		}),
	)
	g.GenerateModelAs(
		"post",
		"PostM",
		gen.FieldIgnore("placeholder"),
		gen.FieldGORMTag("postID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_post_postID")
			return tag
		}),
	)
	g.GenerateModelAs(
		"casbin_rule",
		"CasbinRuleM",
		gen.FieldRename("ptype", "PType"),
		gen.FieldIgnore("placeholder"),
	)
}
