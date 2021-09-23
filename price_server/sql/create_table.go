package sql

import (
	"github.com/jmoiron/sqlx"
)

var db_create = `
create database if not exists db_price;
`

var t_coin_history_info = `
create table if not exists t_coin_history_info
(
     id bigint(20) not NULL AUTO_INCREMENT primary key,
     symbol varchar(20) not null,
     timestamp integer not null,
     price decimal(28,8) not null,
     price_origin varchar(20) not null,
	 weight integer not null
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

var t_log_info = `
create table if not exists t_log_info
(
     id bigint(20) not NULL AUTO_INCREMENT primary key,
     client_ip varchar(20) not null,
     method varchar(20) not null,
     post_data varchar(2048),
     request_proto varchar(20) not null,
	 request_time varchar(64) not null,
	 user_agent varchar(256) not null,
	 request_url varchar(64) not null,
	 response_time varchar(64) not null,
	 request_response varchar(2048) not null
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

func createOrderTables(db *sqlx.DB, dbName string) {
	db.MustExec(db_create)
	db.MustExec("USE " + dbName)
	db.MustExec(t_coin_history_info)
	db.MustExec(t_log_info)
}
