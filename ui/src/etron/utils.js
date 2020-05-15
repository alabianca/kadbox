

const sleep = (time = 1000) => {
    return new Promise((res) => {
        setTimeout(() => {
            res()
        }, time)
    })
}

module.exports = {
    sleep,
}