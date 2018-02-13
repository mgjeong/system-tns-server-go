/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/
package main

    import (
        _ "github.com/go-sql-driver/mysql"
        "database/sql"
        "fmt"
    )

    func main() {
        db, err := sql.Open("mysql", "root:password@/tnsDB?charset=utf8")
        checkErr(err)

        // insert
        stmt, err := db.Prepare("INSERT TNSserver SET topic=?,endpoint=?,service_id=?")
        checkErr(err)

        res, err := stmt.Exec("/SEVT/1F3C/N12/visiontester112/", "112.1.12.143:8455", "device-service112[].k{")
        checkErr(err)

        id, err := res.LastInsertId()
        checkErr(err)

        fmt.Println(id)
        // update
        stmt, err = db.Prepare("update TNSserver set topic=? where uid=?")
        checkErr(err)

        res, err = stmt.Exec("/SEVT/2F43G/N11/raw-datas/", id)
        checkErr(err)

        affect, err := res.RowsAffected()
        checkErr(err)

        fmt.Println(affect)

        // query
        rows, err := db.Query("SELECT * FROM TNSserver")
        checkErr(err)

        for rows.Next() {
            var uid int
            var topic string
            var endpoint string
            var service_id string
            err = rows.Scan(&uid, &topic, &endpoint, &service_id)
            checkErr(err)
            fmt.Println(uid)
            fmt.Println(topic)
            fmt.Println(endpoint)
            fmt.Println(service_id)
        }

        // delete
        stmt, err = db.Prepare("delete from TNSserver where uid=?")
        checkErr(err)

        res, err = stmt.Exec(id)
        checkErr(err)

        affect, err = res.RowsAffected()
        checkErr(err)

        fmt.Println(affect)

        db.Close()

    }

    func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }
