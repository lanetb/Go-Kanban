htmx.onLoad(function (content) {   
    const deleteToggleButton = document.getElementById('delete-toggle')
    const deleteButton = content.querySelectorAll('[data-delete-button]')
    const projects = content.querySelectorAll('[data-project]')

    deleteToggleButton.addEventListener('click', () => {
        deleteButton.forEach(button => {
            if (button.classList.contains('delete-notactive')) {
                button.classList.remove('delete-notactive')
            } else {
                button.classList.add('delete-notactive')
            }
        });
        projects.forEach(project => {
            if (project.classList.contains('deleteable')) {
                project.classList.remove('deleteable')
            } else {
                project.classList.add('deleteable')
            }
        });
    });
});

