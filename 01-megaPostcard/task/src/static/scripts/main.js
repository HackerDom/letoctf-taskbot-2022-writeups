const button = document.getElementById('submit-button')
const input = document.querySelector('input[id="card-text"]');
const textErrorSpan = document.querySelector('span[class="caption"]');
const backgroundErrorSpan = document.querySelector('span[class="backgrounds-caption"]');
const radioButtons = document.getElementsByName("card-background__item")

for (const radio of radioButtons) {
    radio.addEventListener(
        "click", radioCleaner
    )
}

input.addEventListener(
    'input', textValidator
)

function radioCleaner (e) {
    backgroundErrorSpan.style.visibility = "hidden";
}
function textValidator (e) {
    input.style.border = "none";
    textErrorSpan.style.visibility = "hidden";
}


button.onclick = function(event) {
    const radio = document.querySelector('input[type=radio]:checked');
    if (input.value === "") {
        input.style.border = "2px solid red";
        textErrorSpan.textContent = "Введите текст11!!11";
        textErrorSpan.style.visibility = "visible";
        event.preventDefault()
        return
    }
    if (input.value.includes("{{")) {
        textErrorSpan.textContent = "Наш сервис не поддерживает шаблоны :("
        textErrorSpan.style.visibility = "visible";
        event.preventDefault()
        return
    }
    if (radio === null) {
        backgroundErrorSpan.textContent = "Нужно выбрать картонку!11"
        backgroundErrorSpan.style.visibility = "visible";
        event.preventDefault()
    }
}

