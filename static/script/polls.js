document.getElementById("addPoll").addEventListener("click", function(){
    var answers = document.getElementById("answers");
    let i = answers.children.length-3;
    var newInput = document.createElement("input");
    newInput.type = "text";
    newInput.name = i;
    newInput.id = i;
    newInput.placeholder = "Vraag "+(i+1);
    answers.appendChild(newInput);
});
document.getElementById("delPoll").addEventListener("click", function(){
    var answers = document.getElementById("answers");
    if (answers.children.length < 5) return;
    answers.removeChild(answers.lastChild);
});
