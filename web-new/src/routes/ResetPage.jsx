import ResetForm from "../schemas/Reset";
import { authService } from "../services";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { toast } from "../scripts/toast.js"

function ResetPage() {
		const [params, _] = useSearchParams();
		const navigate = useNavigate();
		
		const reset = async (data) => {
                let success = true
				try {
						data.token = params.token;
						await authService.resetPassword(data);
				} catch (e) {
                        success = false
						toast(e.message);
				}
                
                if (success) {
                    toast("Salasana päivitetty onnistuneesti")
                    navigate("/login", { replace: true });
                }
                    
		}
		
		return (
				<>
            <h3>Nollaa salasana</h3>
						<ResetForm
								onSubmit={reset}
						/>
				</>
		);
}

export default ResetPage;