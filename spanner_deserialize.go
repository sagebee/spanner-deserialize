//usr/bin/env go run "$0" "$@"; exit "$?"
package main 

	/*
	 * First:
	 * 
	 * Set up emulator & disable IAM ~
	 * 
	 * $ gcloud config configurations create emulator
	 * $ gcloud config set auth/disable_credentials true
	 * $ gcloud config set project test-project 
	 * $ gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
	 * 
	 * Start emulator ~ 
	 * $ gcloud beta emulators spanner start
	 * 
	 * In working term: 
	 * $ export SPANNER_EMULATOR_HOST=0.0.0.0:9010
	 * Set up instance ~ 
	 * gcloud spanner instances create test-instance --config=emulator-config --description=”test-instance” --nodes=1
	*/


	/* Sanity checks:
	 * instances: gcloud spanner instances list
	 * project id: gcloud config list --format 'value(core.project)'
	*/

import(
	"fmt"
	"log"
	"os"
	"strings"
	"context"
	"io/ioutil"
	"regexp"
	//"reflect"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

// TODO flags
// ~ change instance name 
// ~ set database na,a
// quickstart scripts

func main(){

	// defaults 
	var instancename = "test-instance";
	var projectname  = "test-project";
	var dbname = "example-db"

	// TODO flags to change names 

	// get ddl file from args 
	argc := os.Args[1:];
	if (len(argc) < 1){
		log.Fatal("Usage: ./spanner_deserialize.go <ddl_file>");
	}
	filename := argc[0];

	// where am i
	wdir,err := os.Getwd();
	if (err != nil) {log.Fatal(err)};

	// exit if file doesn't exist
	if _, err := os.Stat(wdir+"/"+filename); os.IsNotExist(err) { 
		log.Fatal(err);
	}

	// read file
	filebytes, err := ioutil.ReadFile(filename)
	if err != nil { log.Fatal(err) }
	filetext := string(filebytes)

	// Seperate DDL & DML
	var ddl []string
	var dml []string
	statements := strings.Split(filetext, ";")
	for _,state := range statements {
		state = strings.TrimSpace(state)
		matchdml, merr :=  regexp.MatchString(`(?is)^\n*\s*(INSERT|UPDATE|DELETE)\s+.+$`, state) 
		if (merr != nil){log.Fatal(merr)}
		if (matchdml){ dml = append(dml, state) } // dml
		if (! matchdml && len(state) > 0){ ddl = append(ddl,state) } // ddl
	}

	// set up db name

	// remove extension
	toks := strings.Split(filename, ".")
	dbname = toks[0]

	// remove folders
	toks = strings.Split(dbname, "/")
	dbname = toks[len(toks)-1];
	fmt.Println("XX"+dbname+"XX") // debug

	return;
	
	// create database 
	createDatabase(dbname, instancename,projectname,ddl)

	// populate database	
	populateDatabase(dbname,instancename,projectname,dml) 

}

/**
 * adapted from https://github.com/GoogleCloudPlatform/golang-samples/
*/
func createDatabase(dbname, instancename,projectname string ,ddl []string){

	// set up client
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancel()

	ctx := context.Background();
	adminCli, err := database.NewDatabaseAdminClient(ctx);
	if err != nil { log.Fatal(err); }

	// parent
	var parentstr = "projects/"+projectname+"/instances/"+instancename;

	
	// create db, read in file
	op, err := adminCli.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          parentstr,
		CreateStatement: "CREATE DATABASE `" + dbname + "`",
		ExtraStatements: ddl,
	})
	if err != nil {
			log.Fatal(err)
	}
	if _, err := op.Wait(ctx); err != nil {
			log.Fatal(err)
	}
	fmt.Println( "Created database +", dbname)
	
}

/**
 * adapted from https://github.com/GoogleCloudPlatform/golang-samples/
*/
func populateDatabase(dbname, instancename, projectname string, dml []string){

	// open client
	var db = "projects/"+projectname+"/instances/"+instancename+"/databases/"+dbname;
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
			log.Fatal(err)
	}
	defer client.Close()

	// set up statements 
	var states []spanner.Statement
	for _,line := range dml {
		states = append(states, spanner.NewStatement(line))
	}
	// debug
	//fmt.Print(states)


	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			stmts := states
			rowCounts, err := txn.BatchUpdate(ctx, stmts)
			if err != nil {
					return err
			}
			fmt.Printf("Executed %d SQL statements using Batch DML.\n", len(rowCounts))
			return nil
	})
	if (err != nil) { log.Fatal(err) }
}





