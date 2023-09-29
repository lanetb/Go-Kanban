htmx.onLoad(function (content) {
    const deleteTaskButtontoggle = document.querySelector('[data-delete-task-toggle]')
    const deleteTaskButtons = content.querySelectorAll('[data-task-delete-btn]')
    const Tasks = content.querySelectorAll('[data-task]')

    deleteTaskButtontoggle.forEach(toggle => {
        toggle.addEventListener('click', () => {
            console.log('clicked')
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
});