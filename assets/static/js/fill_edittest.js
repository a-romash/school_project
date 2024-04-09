var url = window.location.href;

// Создаем объект URL
var urlObj = new URL(url);

// Получаем параметр "t" из URL
var parameterT = urlObj.searchParams.get('test_id');

if (parameterT !== '' && parameterT !== null) {
    const options = {
        method: 'POST',
        headers: {
        'Content-Type': 'application/json'
        // Если требуется, можно также добавить другие заголовки
        },
    };
    fetch('/api/v1/gettest?t='+parameterT, options).then(function(response) {
        if (!response.ok) {
          throw new Error('Произошла ошибка при запросе: ' + response.status);
        }
        return response.json();
      }).then(function(response) {
        document.querySelector('#write_title').value = response.title
        response.questions.forEach((question, index) => {
            if (index !== 0) {
                document.querySelector('#add_questions').click()
            }
            var fieldset = document.querySelectorAll('.fs_questions')[index + 1]
            fieldset.querySelector('#question').value = question.question

            var add_answ = fieldset.querySelector('#add_answ')
            var switch_type = fieldset.querySelector('#type_btn')

            if (question.answ_opt === undefined) {
                fieldset.querySelector('#answer').value = response.answers[index]
            } else {
                switch_type.click()

                for (i=1; i<question.answ_opt.length; i++) {
                    add_answ.click()
                }

                var inp_answs = fieldset.querySelectorAll('.inp_answ')
                var rb_answs = fieldset.querySelectorAll('.rb_ans')
                question.answ_opt.forEach((answer, answ_index) => {
                    inp_answs[answ_index].value = answer
                    if (answ_index.toString() === response.answers[index]) {
                        rb_answs[answ_index].checked = true
                    }
                })
            }

        })
      }).catch(function(error) {
        console.error('Сетевая ошибка: ', error);
      });
}