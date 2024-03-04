var button = document.getElementById("enter");
var error = document.getElementById("error")

// Добавление слушателя событий для кнопки
button.addEventListener("click", submitForm);

function submitForm() {
    var login = document.getElementById('login').value;
    var pw = document.getElementById('pw').value;

    var data = {
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
                console.log('Data sent successfully');
                window.location.href = "/";
            } else {
                error.style.display = 'block';
                console.log('Got error');
            }
        }
    }
}