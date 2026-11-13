import LoginForm from "../schemas/LoginForm";
import { createResource } from "solid-js";
import { authService } from "../services";

function LoginPage() {
		const handleLogin = async (data) => {
				try {
						await authService.login(data);
						await refetch();
						navigate("/", { replace: true });
				} catch (e) {
						alert(e.message);
						console.error(e.message);
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