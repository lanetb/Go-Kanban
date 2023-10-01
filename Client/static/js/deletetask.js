htmx.onLoad(function (content) {
    const deleteTaskButtontoggle = document.getElementById('delete-task-toggle')
    const deleteTaskButtons = content.querySelectorAll('[data-task-delete-btn]')
    const Tasks = content.querySelectorAll('[data-task]')

    (deleteTaskButtons.parentNode.parentNode.parentNode)
    deleteTaskButtontoggle.addEventListener('click', () => {
        deleteTaskButtons.forEach(button => {
            if (button.classList.contains('delete-notactive')) {
               button.classList.remove('delete-notactive')
            } else {
                button.classList.add('delete-notactive')
            }
        });
        Tasks.forEach(task => {
            if (task.classList.contains('deleteable')) {
                task.classList.remove('deleteable')
            } else {
                task.classList.add('deleteable')
            }
        });
    });
});