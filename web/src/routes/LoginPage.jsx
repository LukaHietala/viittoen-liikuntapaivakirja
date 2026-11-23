import LoginForm from "../schemas/LoginForm";
import { createResource } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { useAuth } from "../context/AuthContext";
import { authService } from "../services";
import { toast } from "../scripts/toast.js";

function LoginPage() {
		const navigate = useNavigate();
		const { refetch } = useAuth();
		
		const handleLogin = async (data) => {
				try {
						await authService.login(data);
						await refetch();
						navigate("/", { replace: true });
				} catch (e) {					
						toast(e.message);
				}
		};
		
		return (
				<>
						<h3>Kirjaudu</h3>
						<LoginForm onSubmit={handleLogin}/>
				</>
		);
}

export default LoginPage;