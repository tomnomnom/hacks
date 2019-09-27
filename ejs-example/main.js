const express = require('express')
const path = require('path')
const session = require('express-session')
const app = express()
const port = 3000

// serve static files from the static dir
app.use('/static', express.static('static'))

// use the express session handler
app.use(session({
    secret: 'idkfa'
}))

// use ejs files in ./views
app.set('views', path.join(__dirname, 'views'));

// use ejs as the default way to display pages
app.set('view engine', 'ejs');

app.get('/logout', (req, res) => {
    req.session.views = 0
    res.redirect(301, '/')
})

app.get('/', (req, res) => {
    if (!req.session.views){
        req.session.views = 0
    }
    req.session.views++

    var people = ['geddy', 'neil', 'alex']

    // render the people template
    res.render('people.ejs', {
        people: people,
        views: req.session.views
    })
})

app.listen(port, () => console.log(`Example app listening on port ${port}!`))

