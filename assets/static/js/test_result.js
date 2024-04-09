const data = JSON.parse(localStorage.getItem("solution_results"));

if (data === null || data === '') {
    localStorage.clear()
    window.location.href = '/'
}

document.querySelector('#result').textContent = data.cur_score + '/' + data.max_score

document.querySelector('#btn_restart').addEventListener('click', function() {
    window.location.href = '/test?t=' + data.test_id
});

localStorage.clear()