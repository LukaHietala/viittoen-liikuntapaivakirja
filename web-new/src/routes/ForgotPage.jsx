import ForgotForm from "../schemas/ForgotForm";
import { authService } from "../services";
import { createResource, createSignal } from "solid-js";
import { toast } from "../scripts/toast.js"

function ForgotPage() {
		const [sent, setSent] = createSignal(false);

		const sendEmail = async (data) => {
            let success = true
                try {
                    await authService.forgotPassword(data);
                } catch (e) {
                    success = false
                    toast(e.message) 
                }
            if (success)
                setSent(true)
		}
		
		return (
				<>
						<Show when={sent()}
									fallback={
											<>
													<h3>Salasanan nollaus</h3>
													<ForgotForm
															onSubmit={sendEmail}
													/>
											</>
									}>
                                <div class="block">
                                    <span>Nollauslinkin sisältävä sähköposti lähetetty</span>
                                </div>
						</Show>
				</>
		);
}

export default ForgotPage;