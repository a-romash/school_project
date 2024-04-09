document.querySelector('#add_questions').addEventListener('click', function() {
    const fieldsetContainer = document.getElementById('fs_container');

    const fs_template = fieldsetContainer.querySelector('.fs_questions')
    const clonedFieldset = fs_template.cloneNode(true);

    clonedFieldset.removeAttribute("id"); // Удаление id для избежания повторов
    clonedFieldset.style.display = "block";

    clonedFieldset.querySelector('#delete_btn').addEventListener('click', function() {
      var element = clonedFieldset;
      if (element && document.querySelectorAll('.fs_questions').length > 2) {
          element.parentNode.removeChild(element);
      } else {
          console.error("Элемент не найден.");
      }
    });

    const btn_type = clonedFieldset.querySelector('.flag')

    btn_type.addEventListener('click', function() {
      const fs = btn_type.parentNode
      if (fs.querySelector('#type_answer_one_input').style.display === 'none') {
        fs.querySelector('#type_answer_one_input').style.display = 'block';
        fs.querySelector('#type_answer_mult_inputs').style.display = 'none';
        btn_type.querySelector('#type_btn').innerHTML = 'Несколько<br/>ответов';
      } else {
        fs.querySelector('#type_answer_one_input').style.display = 'none';
        fs.querySelector('#type_answer_mult_inputs').style.display = 'block';
        btn_type.querySelector('#type_btn').innerHTML = 'Один<br/>ответ';
      }
    });


    const btn_add = clonedFieldset.querySelector('#add_answ')

    btn_add.addEventListener('click', function() {
      var templ = `
      <li class="answ">
        <input type="radio" class="rb_ans" name="options` + document.querySelectorAll('.fs_questions').length.toString() + `"/>
        <input class="inp_answ" type="text" placeholder="Введите ответ на вопрос">
      </li>
      `;

      btn_add.parentNode.querySelector('#answ_opts').innerHTML = templ
    })
    fieldsetContainer.appendChild(clonedFieldset);

    btn_add.click()
  });

  document.querySelector('#add_questions').click()
