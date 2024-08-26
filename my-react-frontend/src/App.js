import React, { useState, useEffect } from 'react';
import './App.css';

const FrontendProyecto = () => {
    const [input, setInput] = useState('');
    const [output, setOutput] = useState('# sistema de archivos EXT2.\n\tDominic Juan Pablo Rueno Perez\n\t202200075');
    const [inputLines, setInputLines] = useState(1);
    const [outputLines, setOutputLines] = useState(3); // Inicialmente hay 3 líneas en el output por defecto

    const handleExecute = async () => {
        // Mostrar mensaje de procesamiento
        setOutput('Procesando la entrada...');
        setOutputLines(1);

        try {
            const response = await fetch('http://localhost:8080/api/submit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ text: input })
            });

            if (!response.ok) {
                throw new Error(`Error del servidor: ${response.status} ${response.statusText}`);
            }

            const data = await response.json();
            setOutput(data.text);
            setOutputLines(data.text.split('\n').length);
        } catch (error) {
            console.error('Error al procesar la solicitud:', error);
            setOutput(`Error al procesar la solicitud: ${error.message}`);
            setOutputLines(1);
        }
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

    useEffect(() => {
        setInputLines(input.split('\n').length || 1);
    }, [input]);

    useEffect(() => {
        setOutputLines(output.split('\n').length || 1);
    }, [output]);

    return (
        <div className="container">
            <div className="button-group">
                <button 
                    id="executeButton" 
                    onClick={() => {
                        handleExecute(); 
                        window.scrollTo({
                            top: 0,
                            behavior: 'smooth' // Desplazamiento suave hacia la parte superior
                        });
                    }}
                >
                    Ejecutar
                </button>
                <label htmlFor="fileInput" className="custom-file-upload">
                    Subir Archivo
                </label>
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
                            id="output"
                            value={output}
                            readOnly
                            placeholder=""
                        ></textarea>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default FrontendProyecto;
