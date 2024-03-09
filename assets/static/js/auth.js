var button = document.getElementById("enter");
var error = document.getElementById("error")
var error_fields = document.getElementById("error_fields")

// Добавление слушателя событий для кнопки
button.addEventListener("click", submitForm);

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

function setCookie(key, value, d) {
    var expires = "";
    if (d) {
        var date = new Date(d);
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = key + "=" + value + expires + "; path=/";
}

function deleteCookie(key) {
    document.cookie = key + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

let t = getCookie("t");
if (t !== "" && t !== null) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/login', true);
    xhr.setRequestHeader('t', t);

    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                window.location.href = "/";
            } else {
                deleteCookie('t')
            }
        }
    }

    xhr.send();
}

function submitForm() {
    var login = document.getElementById('login').value;
    var pw = document.getElementById('pw').value;

    error_fields.style.display = 'none';
    error.style.display = 'none';

    if (login === '' || pw === '') {
        error_fields.style.display = 'block';
        return
    }

    let data = {
        login: login,
        pw: pw
    };

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/login', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(data));

    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                var response = JSON.parse(xhr.responseText);
                setCookie("t", response.t, response.expires_at)
                window.location.href = "/";
            } else {
                error.style.display = 'block';
                console.log('Got error');
            }
        }
    }
}