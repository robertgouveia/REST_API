import { useNavigate, useParams } from 'react-router-dom'
import { API_URL } from './App'

export const ConfirmationPage = () => {
  const { token } = useParams()
  const redirect = useNavigate()
  const handleConfirm = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: 'PUT',
    })

    if (response.ok) {
      // redirect homepage
      redirect('/')
    } else {
      alert('Failed to confirm')
    }
  }

  return (
    <div>
      <h2>Confirmation</h2>
      <button onClick={handleConfirm}>Click to confirm</button>
    </div>
  )
}