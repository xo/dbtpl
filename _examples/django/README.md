# About django examples

The `django` example is the result of running `dbtpl` against all the supported
databases that Django and `dbtpl` supports, with Django models similar to the
`dbtpl` booktest schema.

## Setup

Install packages:

```sh
# install mysql, postgres, sqlite3 dependencies
$ sudo apt install libpq-dev libmysqlclient-dev libsqlite3-dev

# install sqlserver dependenices
# manually add the microsoft-prod ppa -- see: https://docs.microsoft.com/en-us/sql/connect/odbc/linux-mac/installing-the-microsoft-odbc-driver-for-sql-server?view=sql-server-ver15
$ sudo apt install unixodbc unixodbc-dev odbcinst msodbcsql18

# aur dependencies:
$ yay -S python-pip python-pipenv unixodbc msodbcsql oracle-instantclient-sdk oracle-instantclient-sqlplus oracle-instantclient-tools

# ensure odbcinst.ini has the relevant sqlserver entry
$ cat /etc/odbcinst.ini
[ODBC Driver 18 for SQL Server]
Description=Microsoft ODBC Driver 18 for SQL Server
Driver=/opt/microsoft/msodbcsql18/lib64/libmsodbcsql-18.0.so.1.1
UsageCount=1

# install oracle dependencies
$ cd /path/to/usql/contrib/godror
$ sudo ./grab-instantclient.sh

# fix oob issue with oracle driver
$ cd /path/to/usql/contrib/godror
$ ./fix-oob-config.sh

# install pipenv
$ pip install --user pipenv

# install packages
$ pipenv install

# update packages
$ pipenv update
```
