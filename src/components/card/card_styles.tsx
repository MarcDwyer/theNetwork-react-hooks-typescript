import styled from 'styled-components'


export const Visit = styled.a`
    position: absolute;
    top: 5px;
    left: 5px;
    width: 85px;
    background-color: #25252E;
    border: none;
    cursor: pointer;
    padding: 5px 5px;
    text-align: center;
    text-decoration: none;

    &:focus {
        outline: 0;
    }
`

export const VisitSpan = styled.span`
    color: #eee !important;
`