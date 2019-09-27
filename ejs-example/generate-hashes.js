const bcrypt = require('bcrypt')
const saltRounds = 9;


(async () => {
    let one = await bcrypt.hash("password123", saltRounds)
    let two = await bcrypt.hash("Qwerty123", saltRounds)
    console.log(one)
    console.log(two)
})()



