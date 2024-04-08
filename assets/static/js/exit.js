function deleteCookie(key) {
    document.cookie = key + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

function getCookie(name) {
    let cookieArr = document.cookie.split(";");

    for (let i = 0; i < cookieArr.length; i++) {
        let cookiePair = cookieArr[i].split("=");
        if (name === cookiePair[0].trim()) {
        return decodeURIComponent(cookiePair[1]);
        }
    }

    return null;
}

var exitButton = document.getElementById("exit");
exitButton.addEventListener("click", function() {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/v1/deletetoken', true)
    xhr.setRequestHeader('Content-Type', 'application/json');

    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status !== 200) {
                console.error('Ошибка при запросе: ' + xhr.status);
            }
        }
        deleteCookie('t')
        window.location.href = '/login'
    };
    const data = {token: getCookie("t")}
    xhr.send(JSON.stringify(data));
});
