const express = require('express')
const path = require('path')
const app = express()
const port = 3000

// use ejs files in ./views
app.set('views', path.join(__dirname, 'views'));

// use ejs as the default way to display pages
app.set('view engine', 'ejs');

app.get('/', (req, res) => {
    var people = ['geddy', 'neil', 'alex']

    // render the people template
    res.render('people.ejs', {people: people})
})

app.listen(port, () => console.log(`Example app listening on port ${port}!`))

