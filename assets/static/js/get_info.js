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

const data = {token: getCookie("t")}

const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
      // Если требуется, можно также добавить другие заголовки
    },
    body: JSON.stringify(data)
  };

fetch('/api/v1/getinfo', options).then(function(response) {
    if (!response.ok) {
      throw new Error('Произошла ошибка при запросе: ' + response.status);
    }
    return response.json();
  }).then(function(fieldsetData) {
    const labelAuthor = document.getElementById('name')
    labelAuthor.append(fieldsetData.name)

    const fieldsetContainer = document.querySelector(".fs_container");
    const templateTestTemplate = document.querySelector("#idktoo");

    for (const testId in fieldsetData.tests) {
      if (fieldsetData.tests.hasOwnProperty(testId)) {
        const clonedFieldset = templateTestTemplate.cloneNode(true);
        clonedFieldset.id = testId
        clonedFieldset.style.display = "block";
        
        const test_name = clonedFieldset.querySelector("#test_name");
        test_name.textContent = fieldsetData.tests[testId].test_name;

        const test_id = clonedFieldset.querySelector("#lbl_id");
        test_id.textContent = 'id='+testId;
        test_id.addEventListener('click', function() {
          navigator.clipboard.writeText(location.host + '/test?t=' + testId)
          alert("Ссылка на тест скопирована в буфер обмена")
      })

        const table_result = clonedFieldset.querySelector("#table_result")

        if (fieldsetData.tests[testId].solutions.length === 0) {
          table_result.append('Решений нет.')
        }

        fieldsetData.tests[testId].solutions.forEach(solution => {
        const solDecodedString = decodeURIComponent(escape(atob(solution)));
        const solDecoded = JSON.parse(solDecodedString);
        
        let tr = document.createElement("tr");

// Создание и заполнение первой ячейки <td>
let td = document.createElement("td");
td.textContent = solDecoded.author;
tr.appendChild(td);

// Создание и заполнение второй ячейки <td>
let td1 = document.createElement("td");
td1.textContent = solDecoded.class;
tr.appendChild(td1);

// Создание и заполнение третьей ячейки <td>
let td2 = document.createElement("td");
td2.textContent = solDecoded.result + '/' + fieldsetData.tests[testId].max_score;
tr.appendChild(td2);

// Создание и заполнение четвертой ячейки <td> с текстовым значением solDecoded.grade
let td3 = document.createElement("td");
td3.textContent = fieldsetData.tests[testId].grade;
console.log(solDecoded.grade)
if (td3.textContent == ""){
  tr.append("undefiend");
} else {
  tr.appendChild(td3);
}

// Добавление строки <tr> к таблице
table_result.appendChild(tr);



        })


        const lbl_stat = clonedFieldset.querySelector("#lbl_stat");
        lbl_stat.append(fieldsetData.tests[testId].amount + " чел.")


        var del_btn = clonedFieldset.querySelector('#delete_btn')
        del_btn.addEventListener('click', function() {
          var xhr = new XMLHttpRequest();
          xhr.open('POST', '/api/v1/deletetest', true)
          xhr.setRequestHeader('Content-Type', 'application/json');
          const id = del_btn.parentNode.parentNode.id
          xhr.onreadystatechange = function() {
              if (xhr.readyState === XMLHttpRequest.DONE) {
                  if (xhr.status !== 200) {
                      console.error('Ошибка при запросе: ' + xhr.status);
                  }
                  if (xhr.status !== 500) {
                    var element = clonedFieldset;
                    if (element) {
                        element.parentNode.removeChild(element);
                    } else {
                        console.error("Элемент не найден.");
                    }
                  }
              }
          };

          const data = {token: getCookie("t"), test_id: id}
          xhr.send(JSON.stringify(data));
      
      });

      var edit_btn = clonedFieldset.querySelector('#edit_btn')
      edit_btn.addEventListener('click', function() {
        window.location.href = '/edit_test?mode=edit&test_id=' + testId
      })
        
        fieldsetContainer.appendChild(clonedFieldset);
      }
  }
  }).catch(function(error) {
    console.error('Сетевая ошибка: ', error);
  });