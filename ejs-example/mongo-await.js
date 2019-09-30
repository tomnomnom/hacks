const MongoClient = require('mongodb').MongoClient
const assert = require('assert')

const url = "mongodb://localhost:27017"
const dbName = "rushwebapp"

const client = new MongoClient(url, {
    useNewUrlParser: true,
    useUnifiedTopology: true
});

(async () => {
    try {
        await client.connect()
        const db = client.db(dbName)
        const users = db.collection('users')

        await users.createIndex({username: 1}, {unique: true})

        await users.insertOne({username: "sam"})

        const docs = await users.find({}).toArray()

        console.log(docs)
    } catch (err) {
        console.log("err", err)
    }
})()

    


