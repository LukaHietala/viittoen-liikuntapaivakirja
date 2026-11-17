import { authService } from "../services";
import { useNavigate } from "@solidjs/router";
import { useAuth } from "../context/AuthContext.jsx";

function LogoutPage() {
		const navigate = useNavigate();
		const { refetch } = useAuth();

		const handleLogout = async () => {
				try {
						await authService.logout();
				} catch (e) {
						console.error(e.message);
				} finally {
						await refetch();
						navigate("/login", { replace: true });
				}
		};

		return (
				<>
						<p>Logout?</p>
						<button onClick={handleLogout}>Yes</button>
				</>
		);
}

export default LogoutPage;