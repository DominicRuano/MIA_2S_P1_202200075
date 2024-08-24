import React, { useState, useEffect } from 'react';
import './App.css';

const FrontendProyecto = () => {
    const [input, setInput] = useState('');
    const [output, setOutput] = useState('# Consola de salida\n\nAquí se mostrará la salida del programa.');
    const [inputLines, setInputLines] = useState(1);
    const [outputLines, setOutputLines] = useState(1);

    const handleExecute = () => {
        setOutput(`Procesando la entrada...`);
        setOutputLines(1);

        // Simular un proceso de ejecución
        setTimeout(() => {
            setOutput(`${input}`);
            setOutputLines(input.split('\n').length);
        }, 2000); // Simulación de 2 segundos
    };

    const handleFileUpload = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                setInput(e.target.result);
            };
            reader.readAsText(file);
        }
    };

    const handleSendText = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/submit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ text: input })
            });
    
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
    
            const data = await response.json();
        } catch (error) {
            console.error('Failed to fetch:', error);
        }
    };
    

    useEffect(() => {
        setInputLines(input.split('\n').length);
    }, [input]);

    useEffect(() => {
        setOutputLines(output.split('\n').length);
    }, []);

    return (
        <div className="container">
            <h3></h3>
            <div className="button-group">
                <button id="executeButton" onClick={() => {
                    handleExecute(); 
                    window.scrollTo({
                        top: 0,
                        behavior: 'smooth'      // Desplazamiento suave hacia la parte superior
                    });
                    handleSendText();
                }}>
                    Ejecutar
                </button>
                <input 
                    type="file" 
                    id="fileInput" 
                    onChange={handleFileUpload} 
                    accept=".txt"
                />
            </div>
            <div className="horizontal-layout">
                <div className="input-section">
                    <label htmlFor="input">Entrada:</label>
                    <div className="textarea-wrapper">
                        <div className="line-numbers">
                            {Array.from({ length: inputLines }, (_, i) => (
                                <span key={i}>{i + 1}</span>
                            ))}
                        </div>
                        <textarea
                            id="input"
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            placeholder="Escribe aquí..."
                        ></textarea>
                    </div>
                </div>
                <div className="output-section">
                    <label htmlFor="output">Salida:</label>
                    <div className="textarea-wrapper">
                        <div className="line-numbers">
                            {Array.from({ length: outputLines }, (_, i) => (
                                <span key={i}>{i + 1}</span>
                            ))}
                        </div>
                        <textarea 
                            disabled
                            id="output"
                            value={output}
                            onChange={(e) => setInput(e.target.value)}
                            placeholder=""
                        ></textarea>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default FrontendProyecto;
