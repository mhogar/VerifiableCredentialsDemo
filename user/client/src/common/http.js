const axios = require('axios')

export default axios.create({
    baseURL: "http://localhost:8082",
    headers: {
        "Content-type": "application/json"
    }
})
