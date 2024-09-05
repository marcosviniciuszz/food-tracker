document.addEventListener('DOMContentLoaded', () => { 
    getOrders();
} );

function getOrders() {

    clearAllItems()

    fetch('http://localhost:8080/orders') 
        .then(response => response.json())
        .then(data => {
            if (Array.isArray(data)) {
                const pendingList = document.getElementById('pending-list');
                const confirmedList = document.getElementById('confirmed-list');
                const preparationList = document.getElementById('preparation-list');
                const readyList = document.getElementById('ready-list');

                data.forEach(pedido => {
                    const listItem = document.createElement('li');
                    listItem.textContent = `Pedido ID: ${pedido._id}`;

                    let button;
                    switch (pedido.status) {
                        case 'pending':
                            button = document.createElement('button');
                            button.className = 'action-button';
                            button.textContent = 'Confirmar';
                            button.onclick = () => handleAction(pedido.orderId, 'confirm');
                            listItem.appendChild(button);
                            pendingList.appendChild(listItem);
                            break;
                        case 'confirmed':
                            button = document.createElement('button');
                            button.className = 'action-button';
                            button.textContent = 'Iniciar Preparo';
                            button.onclick = () => handleAction(pedido.orderId, 'startPrepare');
                            listItem.appendChild(button);
                            confirmedList.appendChild(listItem);
                            break;
                        case 'preparation':
                            button = document.createElement('button');
                            button.className = 'action-button';
                            button.textContent = 'Saiu pra Entrega';
                            button.onclick = () => handleAction(pedido.orderId, 'dispatch');
                            listItem.appendChild(button);
                            preparationList.appendChild(listItem);
                            break;
                        case 'ready':
                            // Itens prontos não terão botão
                            readyList.appendChild(listItem);
                            break;
                        default:
                            break;
                    }
                });
            }
        })
        .catch(error => {
            console.error('Erro ao carregar pedidos:', error);
        });
}

function clearAllItems() {

    const lists = document.querySelectorAll('#pending-list, #confirmed-list, #preparation-list, #ready-list');
    
    lists.forEach(list => {
        list.innerHTML = ''; 
    });
}


function handleAction(id, action) {
    updateOrder(id, action);
}

function updateOrder(id, action) {
    let url = '';
    switch (action) {
        case 'confirm':
            url = `http://localhost:8080/orders/${id}/confirm`;
            break;
        case 'startPrepare':
            url = `http://localhost:8080/orders/${id}/startPrepare`;
            break;
        case 'dispatch':
            url = `http://localhost:8080/orders/${id}/dispatch`;
            break;
        default:
            console.error('Ação desconhecida:', action);
            return;
    }

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(() => {
            getOrders();
        })
        .catch(error => console.error('Erro:', error));
}

