import UserList from "../components/admin/UserList";
import { createResource, createSignal } from "solid-js";
import UserForm from "../schemas/UserForm";
import { userService } from "../services";
import { useAuth } from "../context/AuthContext.jsx";
import { toast } from "../scripts/toast.js";

function AdminPage() {
		const { user } = useAuth();
		const [adding, setAdding] = createSignal(false);
		
		const [users, { refetch: refetchUsers }] = createResource(
				() => user() && user().id,
				async () => {
						return await userService.listAll();
				},
		);

		const addUser = async (data) => {
				try {
						await userService.create(data);
				} catch (e) {
						toast(e.message);
				}
				setAdding(false);
				toast(`Lisättiin käyttäjä: ${data.name}`)
				await refetchUsers();
		};

		const editUser = async (data, id) => {
				try {
						await userService.updateByID(id, data);
				} catch (e) {
						toast(e.message);
				}
				toast(`Muokattiin käyttäjää: ${data.name}`)
				await refetchUsers();
		};

		const deleteUser = async (id) => {
				try {
						await userService.deleteByID(id);
				} catch (e) {
						toast(e.message);
				}
				toast("Poistettiin käyttäjä")
				await refetchUsers();
		};
		
		return (
				<>
						<h2>Ylläpito</h2>
                <Show when={adding()}
                      fallback={
                          <div class="block">

															<button onClick={() => setAdding(true)}>
																	Lisää käyttäjä
															</button>
                          </div>
                      }>
                    <UserForm
                        onSubmit={addUser}
                        user={null}
                        cancel={setAdding(false)}
                    />
                </Show>
            <h3>Käyttäjät</h3>
                <Suspense fallback={<p>Ladataan...</p>}>
                    <UserList
                        users={users}
                        editUser={editUser}
                        deleteUser={deleteUser}
                    />
                </Suspense>
				</>
		);
}

export default AdminPage;