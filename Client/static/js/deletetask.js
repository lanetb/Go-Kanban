htmx.onLoad(function (content) {
    const deleteTaskButtonToggles = content.querySelectorAll('[data-delete-task-toggle]')

    deleteTaskButtonToggles.forEach(toggle => {
        toggle.addEventListener('click', () => {
            const board = toggle.closest('[data-board]');
            const deleteTaskButtons = board.querySelectorAll('[data-task-delete-btn]');
            const tasks = board.querySelectorAll('[data-task]');

            deleteTaskButtons.forEach(button => {
                button.classList.toggle('delete-notactive', !button.classList.contains('delete-notactive'));
            });

            tasks.forEach(task => {
                task.classList.toggle('deleteable', !task.classList.contains('deleteable'));
            });
        });
    });
});