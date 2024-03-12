// Данные для заполнения fieldset
// Получаем текущий URL страницы
var url = window.location.href;

// Создаем объект URL
var urlObj = new URL(url);

// Получаем параметр "t" из URL
var parameterT = urlObj.searchParams.get('t');

fetch('/api/gettest?t='+parameterT).then(function(response) {
    if (!response.ok) {
      throw new Error('Произошла ошибка при запросе: ' + response.status);
    }
    return response.json();
  }).then(function(fieldsetData) {
    console.log(fieldsetData)
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
                radioButton.name = "options";

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
