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

function deleteCookie(key) {
    document.cookie = key + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

let t = getCookie("t");
if (t !== "" && t !== null) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/v1/login', true);
    xhr.setRequestHeader('t', t);

    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status !== 200) {
                deleteCookie('t')
                window.location.href = "/login";
            }
        }
    }

    xhr.send();
} else if (window.location.pathname !== '/login' && window.location.pathname !== '/register') {
    window.location.href = '/login'
}
