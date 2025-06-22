import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { FormProvider } from './contexts/FormContext';
import UserInput from './pages/UserInput';
import UserConfirm from './pages/UserConfirm';
import UserComplete from './pages/UserComplete';
import './App.css';
import './styles/components.css';
import './styles/validation.css';
import './styles/pages.css';
import './styles/session.css';

function App() {
  return (
    <FormProvider>
      <Router>
        <div className="App">
          <header className="App-header">
            <h1>会員登録</h1>
          </header>
          <main className="App-main">
            <Routes>
              <Route path="/" element={<UserInput />} />
              <Route path="/confirm" element={<UserConfirm />} />
              <Route path="/complete" element={<UserComplete />} />
            </Routes>
          </main>
        </div>
      </Router>
    </FormProvider>
  );
}

export default App;