var button = document.getElementById("submit");

var error_unknown = document.getElementById("error_unknown")
var error_login = document.getElementById("error_login")
var error_pw = document.getElementById("error_pw")
var error_all = document.getElementById("error_all")

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

function submitForm() {
    var login = document.getElementById('login').value;
    var name = document.getElementById('name').value;
    var surname = document.getElementById('surname').value;
    var school = document.getElementById('sch').value;
    var pw = document.getElementById('pw').value;
    var pw_rep = document.getElementById('pw_rep').value;

    error_unknown.style.display = 'none';
    error_login.style.display = 'none';
    error_all.style.display = 'none';
    error_pw.style.display = 'none';

    if (login === '' || name === '' || surname === '' || school === '' || pw === '' || pw_rep === '') {
        error_all.style.display = 'block';
        return
    }
    if (pw !== pw_rep) {
        error_pw.style.display = 'block';
        return
    }

    var data = {
        login: login,
        name: name,
        lastname: surname,
        school: school,
        pw: pw
    };

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/register', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(data));

    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                window.location.href = "/login";
            } else {
                if (xhr.status === 409) {
                    error_login.style.display = 'block';
                } else {
                    error_unknown.style.display = 'block';
                }
            }
        }
    }
}