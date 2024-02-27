mongosh -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" admin <<EOF
db = db.getSiblingDB("$MONGO_INITDB_NAME");
db.createUser({
'user': "$MONGO_NEWUSER_NAME",
'pwd': "$MONGO_NEWUSER_PASSWORD",
'roles': [
      {'role': 'dbOwner', 'db': "$MONGO_INITDB_NAME"}
   ]
});
db.createCollection("$MONGO_INITDB_COL_USER", {
   validator: {
      \$jsonSchema: {
         bsonType: "object",
         required: [ "email" ],
         properties: {
            email: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            refresh_token: {
               bsonType: "string",
               description: "must be a string if the field exist"
            },
            expires_at: {
               bsonType: "long",
               description: "must be a long if the field exist"
            },
         }
      }
   }
});
db.$MONGO_INITDB_COL_USER.createIndex({ email: 1 },{ unique: true });
db.$MONGO_INITDB_COL_USER.createIndex({ refresh_token: 1 }, { unique: true, sparse: true});
EOF