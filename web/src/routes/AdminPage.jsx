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
        let success = true
				try {
						await userService.create(data);
				} catch (e) {
            success = false
						toast(e.message);
				}
        if (!success) return;
        setAdding(false);
        toast(`Lisättiin käyttäjä: ${data.name}`)
				await refetchUsers();        
		};

		const editUser = async (data, id) => {
        let success = true
				try {
						await userService.updateByID(id, data);
				} catch (e) {
            success = false
						toast(e.message);
				}
				if (!success) return;
        toast(`Muokattiin käyttäjää: ${data.name}`)
				await refetchUsers();
		};

		const deleteUser = async (id) => {
        let success = true
				try {
						await userService.deleteByID(id);
				} catch (e) {
            success = false
						toast(e.message);
				}
        if (!success) return;
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
                    title="Luodaan uusi käyttäjä"
                    description="Kun käyttäjä luodaan tälle lähetetään tunnukset
                    sähköpostiin. Jos sähköpostia ei tule, voi yrittää
                    esimerkiksi hakea uutta salasanaa &quot;Unohditko
                    salasana&quot; sivulta."
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