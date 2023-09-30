
htmx.onLoad(function (content) {
    const draggables = content.querySelectorAll('.task-container');
    const droppables = content.querySelectorAll('.tasks');
    const board = content.querySelector('.board-container');

    draggables.forEach((task) => {
        task.addEventListener("dragstart", () => {
            task.classList.add("is-dragging");
        });
        task.addEventListener("dragend", () => {
            console.log(`{ "boardID": "${currentzone.id}", "projectID": "${board.id}", "taskID": "${task.id}" }`);
            htmx.ajax('POST', '/onDragEnd/', {
                swap: 'none',
                values: {
                    boardID: currentzone.id,
                    projectID: board.id,
                    taskID: task.id
                }
            });
            task.classList.remove("is-dragging");
        });
    });

    droppables.forEach((zone) => {
        zone.addEventListener("dragover", (e) => {
            e.preventDefault();

            const bottomTask = insertAboveTask(zone, e.clientY);
            const curTask = content.querySelector(".is-dragging");

            if (!bottomTask) {
                zone.appendChild(curTask);
            } else {
                zone.insertBefore(curTask, bottomTask);
            }
            currentzone = zone;
        });
    });

    const insertAboveTask = (zone, mouseY) => {
        const els = zone.querySelectorAll(".task-container:not(.is-dragging)");

        let closestTask = null;
        let closestOffset = Number.NEGATIVE_INFINITY;

        els.forEach((task) => {
            const { top } = task.getBoundingClientRect();

            const offset = mouseY - top;
            if (offset < 0 && offset > closestOffset) {
                closestOffset = offset;
                closestTask = task;
            }
        });

        return closestTask;
    };

    function parse(str) {
        var args = [].slice.call(arguments, 1),
            i = 0;

        return str.replace(/%s/g, () => args[i++]);
    }
});