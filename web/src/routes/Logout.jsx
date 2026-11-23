import { authService } from "../services";
import { useNavigate } from "@solidjs/router";
import { useAuth } from "../context/AuthContext.jsx";
import { toast } from "../scripts/toast.js";

function LogoutPage() {
		const navigate = useNavigate();
		const { refetch } = useAuth();

		const handleLogout = async () => {
				try {
						await authService.logout();
				} catch (e) {
						toast(e.message);
				} finally {
						await refetch();
						navigate("/login", { replace: true });
				}
		};

		return (
            <>
                <h2>Kirjaudu ulos?</h2>
				<div class="block">
						<button onClick={handleLogout}>Kyllä</button>
				</div>
            </>
		);
}

export default LogoutPage;