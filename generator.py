import MySQLdb
import os

dbconnect = MySQLdb.connect(os.getenv("DB_HOST"), os.getenv("DB_USERNAME"), os.getenv("DB_PASSWORD"), os.getenv("DB_DATABASE"))

cursor = dbconnect.cursor()
cursor.execute("SELECT * FROM information_schema.columns WHERE table_schema = 'expirationdotdev_local';")


def convert_to_go(name):
    names = name.split("_")
    k = 0
    for v in names:
        names[k] = v.capitalize()
        k=k+1
    return "".join(names)


objects = {}

data = cursor.fetchall()
for row in data:
    table_name = convert_to_go(row[2])
    column_name = convert_to_go(row[3])
    nullable = row[6]
    data_type = row[7]
    if objects.get(table_name) is None:
        objects[table_name] = {}
    objects[table_name][column_name] = {}
    objects[table_name][column_name]["nullable"] = nullable == "NO"
    if data_type in ["longtext", "varchar"]:
        objects[table_name][column_name]["data_type"] = "string"

    if data_type in ["bigint", "int"]:
        objects[table_name][column_name]["data_type"] = "int"

    if data_type in ["float"]:
        objects[table_name][column_name]["data_type"] = "float64"

print(f"""package app""")

for key in objects:
    print(f"""
type {key} struct {{""")
    struct = objects[key]
    for attribute in struct:
        if struct[attribute].get("data_type") is None:
            continue
        data_type = struct[attribute]["data_type"]
        print(f"""  {attribute} {data_type}""")

    print(f"""}}""")

dbconnect.close()