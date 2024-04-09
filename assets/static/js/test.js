// Данные для заполнения fieldset
// Получаем текущий URL страницы
var url = window.location.href;

// Создаем объект URL
var urlObj = new URL(url);

// Получаем параметр "t" из URL
var parameterT = urlObj.searchParams.get('t');

fetch('/api/v1/gettest?t='+parameterT).then(function(response) {
    if (!response.ok) {
      throw new Error('Произошла ошибка при запросе: ' + response.status);
    }
    return response.json();
  }).then(function(fieldsetData) {
    const labelAuthor = document.getElementById('author')
    labelAuthor.append(fieldsetData.author)

    const labelTitle = document.getElementById('name')
    labelTitle.append(fieldsetData.title)

    const options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false, timeZone: 'UTC' };
    var date = new Date(fieldsetData.created);
    var formattedDate = date.toLocaleString('ru-RU', options);

    const labelCreated = document.getElementById('created')
    labelCreated.append(formattedDate)

    date = new Date(fieldsetData.updated);
    formattedDate = date.toLocaleString('ru-RU', options);

    const labelUpdated = document.getElementById('updated')
    labelUpdated.append(formattedDate)

    const fieldsetContainer = document.getElementById("fs_container");
    const templateInputFieldset = document.getElementById("fs_question_input");
    const templateRBFieldset = document.getElementById("fs_question_rb");

    // Заполнение и добавление fieldset на основе шаблона
    fieldsetData.questions.forEach(data => {
        if (data.answ_opt === undefined) {
                const clonedFieldset = templateInputFieldset.cloneNode(true);
                clonedFieldset.removeAttribute("id"); // Удаление id для избежания повторов
                clonedFieldset.style.display = "block";
            
                const questionLabel = clonedFieldset.querySelector("#question");
            
                questionLabel.textContent = data.question;
            
                fieldsetContainer.appendChild(clonedFieldset);
        } else {
            const clonedFieldset = templateRBFieldset.cloneNode(true);
            clonedFieldset.removeAttribute("id"); // Удаление id для избежания повторов
            clonedFieldset.style.display = "block";
        
            const questionLabel = clonedFieldset.querySelector("#question");
        
            questionLabel.textContent = data.question;

            var ul = clonedFieldset.querySelector("ul")
            data.answ_opt.forEach(answ => {
                const label = document.createElement("label");
                

                const radioButton = document.createElement("input");
                radioButton.type = "radio";
                radioButton.name = "options" + document.querySelectorAll('#answ_opts').length.toString();

                const div = document.createElement("li");
                label.appendChild(radioButton);
                label.append(" " + answ);
                div.appendChild(label);


                ul.appendChild(div);
            });
        
            fieldsetContainer.appendChild(clonedFieldset);
        }
    });

  }).catch(function(error) {
    console.error('Сетевая ошибка: ', error);
    window.location.href = '/'
  });
