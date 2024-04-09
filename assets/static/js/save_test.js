document.querySelector('.save_button').addEventListener('click', function() {
    var url = window.location.href;

    // Создаем объект URL
    var urlObj = new URL(url);

    // Получаем параметр "t" из URL
    var parameterT = urlObj.searchParams.get('test_id');

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

    if (document.querySelector('#write_title').value === '') {
        alert("Не все поля заполнены!")
    } else {
        var data = {
            token: getCookie("t"),
            title: document.querySelector('#write_title').value,
            questions: [],
            answers: []
        }

        var err = false
            document.querySelectorAll('.fs_questions').forEach(fs_quest => {
                if (err) {
                    return
                }

                if (fs_quest.style.display !== "none") {
                    var title = fs_quest.querySelector("#question")
                    if (title.value === '') {
                        err = true
                        return
                    }
                    var one_inp = fs_quest.querySelector('#type_answer_one_input')
                    var mult_inp = fs_quest.querySelector('#type_answer_mult_inputs')
                    if (mult_inp.style.display === "none") {
                        if (one_inp.querySelector('#answer').value === '') {
                            err = true
                            return
                        }
                        data.questions.push(
                            {
                                "question": title.value
                            }
                        )
                        data.answers.push(one_inp.querySelector('#answer').value)
                    } else {
                        var ans = -1
                        var k = 0
                        var answers = []
                        mult_inp.querySelectorAll('.answ').forEach(answer => {
                            var answer_input = answer.querySelector('.inp_answ')
                            if (answer_input.value === '') {
                                err = true
                                return
                            }

                            answers.push(answer_input.value)
                            if (answer.querySelector('.rb_ans').checked) {
                                ans = k
                            }
                            k++
                        })
                        if (ans === -1) {
                            err = true
                            return
                        }
                        data.questions.push(
                            {
                                "question": title.value,
                                "answ_opt": answers
                            }
                        )

                        data.answers.push(ans.toString())
                    }
                }
            })


        if (err) {
            alert("Не все поля заполнены/выбраны!")
        } else {
            if (parameterT === '' || parameterT === null) {
                var xhr = new XMLHttpRequest();
                xhr.open('POST', '/api/v1/createtest');
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.onreadystatechange = function() {
                    if (xhr.readyState === XMLHttpRequest.DONE) {
                        if (xhr.status === 200) {
                            window.location.href = "/";
                        } else {
                            console.log('Got error');
                            alert("Got error (Status code: " + xhr.status + ")")
                        }
                    }
                }
                xhr.send(JSON.stringify(data));
            } else {
                var xhr = new XMLHttpRequest();
                xhr.open('POST', '/api/v1/updatetest');
                xhr.setRequestHeader('Content-Type', 'application/json');
                xhr.onreadystatechange = function() {
                    if (xhr.readyState === XMLHttpRequest.DONE) {
                        if (xhr.status === 200) {
                            window.location.href = "/";
                        } else {
                            console.log('Got error');
                            alert("Got error (Status code: " + xhr.status + ")")
                        }
                    }
                }
                data.test_id = parameterT
                xhr.send(JSON.stringify(data));
            }
        }
    }
});