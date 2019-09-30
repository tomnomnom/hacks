const MongoClient = require('mongodb').MongoClient
const express = require('express')
const path = require('path')
const session = require('express-session')
const bcrypt = require('bcrypt')
const app = express()
const port = 3000

const url = "mongodb://localhost:27017"
const dbName = "rushwebapp"

const client = new MongoClient(url, {
    useNewUrlParser: true,
    useUnifiedTopology: true
});

// mock user databas
// serve static files from the static dir
app.use('/static', express.static('static'))

// use the express session handler
app.use(session({
    secret: 'idkfa',
    resave: false,
    saveUninitialized: true
}))

// use ejs files in ./views
app.set('views', path.join(__dirname, 'views'));

// use ejs as the default way to display pages
app.set('view engine', 'ejs');

// we're going to POST urlencoded data
app.use(express.urlencoded())

app.get('/register', (req, res) => {
    res.render('register')
})

app.post('/register', async (req, res) => {
    try {
        let saltRounds = 9
        let hash = await bcrypt.hash(req.body.password, saltRounds)

        const db = client.db(dbName)
        const users = db.collection('users')

        await users.createIndex({username: 1}, {unique: true})
        await users.insertOne({
            username: req.body.username,
            hash: hash
        })

        res.redirect('/')
    } catch(err) {
        console.log("registration error: ", err)
        res.render("error", {
            "message": "Failed to rernaegister. Sorry."
        })
    }
})

app.post('/login', async (req, res) => {

    try {
        const db = client.db(dbName)
        const users = db.collection('users')

        const user = await users.findOne({username: req.body.username})

        let success = await bcrypt.compare(req.body.password, user.hash)

        if (!success){
            throw "Bad username or password!"
        }

        req.session.user = req.body.username
        res.redirect(302, '/')

    } catch(err) {
        res.render("error", {
            message: "Bad username or password!"
        }) 
    }
})

app.get('/logout', (req, res) => {
    req.session.user = null
    res.redirect(302, '/')
})

app.get('/', (req, res) => {
    var people = ['geddy', 'neil', 'alex']

    // render the people template
    res.render('people.ejs', {
        people: people,
        user: req.session.user
    })
})

app.get('/api/users', async (req, res) => {

    try {
        const db = client.db(dbName)
        const users = await db.collection('users').find({
            username: { $regex: '^' + req.query.q }
        }).toArray()

        res.json(users.map(u => u.username))

    } catch(err) {
        res.render("error", {
            message: "Failed to list users"
        }) 
    }
})

app.get('/users', (req, res) => {
    res.render('users', {
        user: req.session.user
    })
})

// Connect to mongo, and then start listening
client.connect(() => {
    app.listen(port, () => console.log(`Example app listening on port ${port}!`))
})

