
const baseUrl = 'http://localhost:8080/api/v1'

const post = (url: string, body: any) => {
    url = baseUrl + url;
    return fetch(url, {
        method: 'POST',
        body: JSON.stringify(body),
    }).then((res) => res.json());
}

const API = {
    post,
}

export default API;