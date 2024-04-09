document.querySelector('#end_test').addEventListener('click', function() {
    const fieldsetContainer = document.getElementById("fs_container");

    var data = {};
    var err = false
    var fieldsets = document.querySelectorAll('fieldset');
    data['author'] = document.getElementById('surname_name').value;
    data['class'] = document.getElementById('clas_num').value;
    if (data['author'] === '' || data['class'] === '') {
      alert("Не все поля заполнены!")
      err = true
      return
    }
    // Получаем текущий URL страницы
    var url = window.location.href;

    // Создаем объект URL
    var urlObj = new URL(url);

    // Получаем параметр "t" из URL
    var test_id = urlObj.searchParams.get('t');
    data.test_id = test_id
    data.answers = []
    
    fieldsetContainer.querySelectorAll("fieldset").forEach(fs => {
      if (fs.querySelector('ul') === null) {
        let answer = fs.querySelector('input').value
        if (answer === "") {
          err = true
        }
        data['answers'] = data['answers'].concat(fs.querySelector('input').value)
      } else {
        var got_answer = false
        var idksomeusefulvariable = 0
        fs.querySelector('ul').querySelectorAll('input').forEach(radio => {
          if (radio.checked) {
            data['answers'] = data['answers'].concat(idksomeusefulvariable.toString())
            got_answer = true
          }
          idksomeusefulvariable++
        })
        if (!got_answer) {
          err = true
        }
      }
    })
    if (err) {alert('Не все поля заполнены!');return;}
    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
        // Если требуется, можно также добавить другие заголовки
      },
      body: JSON.stringify(data)
    };
    fetch('/api/v1/getresult', options).then(function(response) {
      if (!response.ok) {
        throw new Error('Произошла ошибка при запросе: ' + response.status);
      }
      return response.json();
    }).then(function(response) {
      localStorage.setItem('solution_results', JSON.stringify(response))
      window.location.href = "/result";
    }).catch(function(error) {
      console.error('Сетевая ошибка: ', error);
    });
    });