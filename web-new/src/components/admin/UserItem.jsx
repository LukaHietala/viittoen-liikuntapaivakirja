import { createSignal } from "solid-js"
import UserForm from "../../schemas/UserForm";

function UserItem(props) {
		const user = props.user;
		
    const [editing, setEditing] = createSignal(false)
		const [deleting, setDeleting] = createSignal(false);

		const editT = (data) => {
				props.editUser(data, user.id);
		};
		
		return (
        <>						
						Käyttäjä: {user.name}
            <Show when={editing()}
									fallback={<button onClick={() => setEditing(true)}>Muokkaa</button>}>
								<UserForm
										user={() => user}
										onSubmit={editT}
										cancel={setEditing(false)}
								/>
            </Show>
						<Show when={deleting()}
									fallback={
											<button onClick={() => setDeleting(true)}>
													Poista
											</button>
									}>
								<button onClick={() => setDeleting(false)}>
										Peru
								</button>
								<button onClick={() => props.deleteUser(user.id)}>
										Varmista
								</button>
						</Show>
        </>
		);
}

export default UserItem;