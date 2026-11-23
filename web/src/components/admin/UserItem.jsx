import { createSignal } from "solid-js";
import UserForm from "../../schemas/UserForm";
import { useAuth } from "../../context/AuthContext";
import ConfirmButton from "../common/ConfirmButton";

function UserItem(props) {
		const { user } = useAuth();
		const cuser = props.user;

		const [editing, setEditing] = createSignal(false);
		const [isConfirmingDelete, setIsConfirmingDelete] = createSignal(false)

		const editT = (data) => {
				props.editUser(data, cuser.id);
		};

		return (
						<Show
								when={!editing()}
								fallback={
										<UserForm
                        title={`Muokataan käyttäjää: ${cuser.name}`}
												user={() => cuser}
												onSubmit={editT}
												cancel={setEditing(false)}
												self={cuser.id === user().id}
										/>
								}
						>
								<div class="item">
										<div class="item-info">
												<span>{cuser.name}</span>
												<Show when={cuser.id === user().id}>
														<span class="badge">Sinä</span>
												</Show>
												<Show when={cuser.isAdmin}>
														<span class="badge">Ylläpitäjä</span>
												</Show>
										</div>
										<div class="item-actions">
												<Show when={!isConfirmingDelete()}>
														<button onClick={() => setEditing(true)}>Muokkaa</button>
												</Show>
												<Show when={user() && user().id != cuser.id}>
														<ConfirmButton
																onConfirmChange={setIsConfirmingDelete} 
																onClick={() => props.deleteUser(cuser.id)}
														/>
												</Show>
										</div>
								</div>
						</Show>
		);
}

export default UserItem;