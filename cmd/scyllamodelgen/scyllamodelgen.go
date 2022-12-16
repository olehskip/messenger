package main

import (
	"log"
	"flag"
	"os"
	"unicode"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

var (
	cmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagCluster = cmd.String("cluster", "127.0.0.1", "ip:port or just ip of the cluster")
	flagKeyspace = cmd.String("keyspace", "", "keyspace to inspect")
	flagTableName = cmd.String("table", "", "table to generate from")
	flagPkgName = cmd.String("package", "models", "the name you wish to assign to your generated package")
	flagModelName = cmd.String("model", "", "the name of your model struct")
	flagOutput = cmd.String("output", "models", "the name of the folder to output to")

	session gocqlx.Session

	packagesToImport = make(map[string]bool)
)

func main() {
	log.SetFlags(0) // disable timestmap and other stuff

	err := cmd.Parse(os.Args[1:])

	if err != nil {
		log.Fatalln("Can't parse flags")	
	}

	if *flagKeyspace == "" {	
		log.Fatalln("Missing required flag: flag keyspace")	
	}
	
	if *flagTableName == "" {	
		log.Fatalln("Missing required flag: flag table")	
	}

	if *flagModelName == "" {
		log.Fatalln("Missing required flag: model name")
	}

	if *flagOutput == "" {
		log.Fatalln("Missing required flag: flag output")
	}

	session, err = connectToDatabase()
	if err != nil {	
		log.Fatalf("Can't connect to database, err = %v", err)	
	} else {
		log.Println("Connected to database")
	}

	columnsDesc, err := getColumnsDesc()
	
	if err != nil {
		log.Fatalln("Can't get information about the table")
	}

	model := getModel(columnsDesc)
	
	writeToFile(model, *flagOutput)

	disconnectFromDatabase()
	
	log.Println("Created model successfully")
}

func connectToDatabase() (gocqlx.Session, error) {
	cluster := gocql.NewCluster(*flagCluster)
	cluster.Keyspace = *flagKeyspace;

	return gocqlx.WrapSession(cluster.CreateSession())
}

func disconnectFromDatabase() {
	session.Close()
}

type ColumnDesc struct {
	ColumnName string
	Type string
}

func (c *ColumnDesc) ToString() string {
	return (convertToCamelCase(c.ColumnName) + " " + matchType(c.Type))
}

func convertToCamelCase(columnName string) string {
	camelCaseName := ""
	
	for i, chr := range columnName {
		if columnName[i] == '_' { 
			continue
		}
		if i == 0 || columnName[i - 1] == '_' {
			chr = unicode.ToUpper(chr)
		}

		camelCaseName += string(chr)
	}
	
	return camelCaseName
}

type SpecialType struct {
	PackageName string
	TypeName string
}

var specialTypeMatching = map[string]SpecialType {
	"timestamp": {PackageName: "time", TypeName: "time.Time"},
}

// Types listed in specialTypeMatching require special care
// We can specify PackagaName for them and it will be automatically added
// If an additional package is not required then PackageName should be = ""
// Types not listed in specialTypematching will be presented as string
func matchType(typeName string) string {
	specialType, ok := specialTypeMatching[typeName]
	if ok {
		if specialType.PackageName != "" {
			packagesToImport[specialType.PackageName] = true
		}

		return specialType.TypeName
	} else {
		return "string"
	}
}

// Get information about columns column name and its type
func getColumnsDesc() ([]ColumnDesc, error) {
	q := session.Query("SELECT column_name, type FROM system_schema.columns WHERE keyspace_name = 'messenger' AND table_name = 'users'", nil)

	var columns []ColumnDesc

	err := q.SelectRelease(&columns)
	return columns, err
}

func getModel(columnsDesc []ColumnDesc) string {
	fields := ""
	for _, columnDesc := range columnsDesc {
		fields += "\t" + columnDesc.ToString() + "\n"
	}

	packagesNames := ""
	for packageName := range packagesToImport {
		packagesNames += "\t" + "\"" + packageName + "\"" + "\n"
	}
	
	res := "package " + *flagPkgName + "\n\n"
	if packagesNames != "" {
		res += "import (\n" + packagesNames + ")\n\n";	
	}
	res += "type " + *flagModelName + " struct {\n" + fields + "}\n"

	return res
}

func writeToFile(model string, path string) {
	file, err := os.Create(path)
    if err != nil {
        log.Fatalln("Can't open file")
    }
    defer file.Close()

    file.WriteString(model)
}

