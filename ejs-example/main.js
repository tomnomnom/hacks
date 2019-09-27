const express = require('express')
const path = require('path')
const session = require('express-session')
const bcrypt = require('bcrypt')
const app = express()
const port = 3000

// mock user database
const users = {
    "sam": {hash: "$2b$09$a20ufU6tABLVLZeZE9Ry9uKRPILZqtmDbLbQQDdHpoJWjKNQYUBz2"},
    "jen": {hash: "$2b$09$DwWboyxs3tz9EubWEH8pouzV54lvleQSNBktUVcy0YmDGF4Efslha"}
}

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

app.post('/login', async (req, res) => {
    if (!users[req.body.username]){
        res.render("error", {
            message: "Bad username or password!"
        }) 
        return
    }

    let user = users[req.body.username]

    let success = await bcrypt.compare(req.body.password, user.hash)

    if (!success){
        res.render("error", {
            message: "Bad username or password!"
        }) 
        return
    }

    req.session.user = req.body.username
    res.redirect(301, '/')
})

app.get('/logout', (req, res) => {
    req.session.user = null
    res.redirect(301, '/')
})

app.get('/', (req, res) => {
    var people = ['geddy', 'neil', 'alex']

    // render the people template
    res.render('people.ejs', {
        people: people,
        user: req.session.user
    })
})

app.listen(port, () => console.log(`Example app listening on port ${port}!`))

